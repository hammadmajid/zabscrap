import type { Actions } from './$types';
import { launch } from '@cloudflare/playwright';
import type { CourseAttendance } from '$lib/models/attendance';

interface ParseTagOptions {
	label: string;
	html: string;
}

/**
 * Parse a value from HTML table row by label
 * Mirrors Go's parseTag function
 */
function parseTag({ html, label }: ParseTagOptions): string {
	const escapedLabel = label.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
	const regex = new RegExp(`(?i)<th[^>]*>${escapedLabel}</th>\\s*<td[^>]*>(.*?)</td>`, 's');
	const match = html.match(regex);
	if (match?.[1]) {
		// Remove HTML tags and trim whitespace
		return match[1].replace(/<[^>]*>/g, '').trim();
	}
	return 'N/A';
}

/**
 * Scrape attendance data using Playwright
 */
async function scrapeAttendance(
	username: string,
	password: string,
	browser: any
): Promise<CourseAttendance[]> {
	const context = await browser.newContext();
	const page = await context.newPage();

	try {
		// Step 1: Navigate to login page and authenticate
		const loginURL = 'https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp';
		await page.goto(loginURL, { waitUntil: 'domcontentloaded' });

		// Fill login form
		await page.fill('input[name="txtLoginName"]', username);
		await page.fill('input[name="txtPassword"]', password);
		await page.fill('input[name="txtCampus_Id"]', '1');

		// Submit form and wait for navigation
		await Promise.all([page.waitForNavigation(), page.click('input[type="submit"]')]);

		// Step 2: Navigate to attendance listing page
		const listURL =
			'https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance';
		await page.goto(listURL, { waitUntil: 'domcontentloaded' });

		// Get page content to extract course parameters
		const pageContent = await page.content();

		// Step 3: Extract course parameters using regex (same as Go version)
		const courseRegex = /chkSubmit\('([^']+)','([^']+)','([^']+)','([^']+)'\)/g;
		const courseMatches = Array.from(pageContent.matchAll(courseRegex));

		const results: CourseAttendance[] = [];

		// Step 4: Fetch each course's attendance data
		for (const match of courseMatches) {
			const [, faculty, semester, section, courseCode] = match;

			// Submit form with course parameters
			await page.goto(listURL, { waitUntil: 'domcontentloaded' });

			// Fill hidden form fields
			await page.evaluate(
				({ faculty, semester, section, courseCode }) => {
					const form = document.querySelector('form') as HTMLFormElement;
					if (form) {
						(form.querySelector('input[name="txtFac"]') as HTMLInputElement).value = faculty;
						(form.querySelector('input[name="txtSem"]') as HTMLInputElement).value = semester;
						(form.querySelector('input[name="txtSec"]') as HTMLInputElement).value = section;
						(form.querySelector('input[name="txtCou"]') as HTMLInputElement).value = courseCode;
					}
				},
				{ faculty, semester, section, courseCode }
			);

			// Submit form and wait for response
			await Promise.all([page.waitForNavigation(), page.click('input[type="submit"]')]);

			// Get course details HTML
			const courseHTML = await page.content();

			// Parse course information
			const courseName = parseTag({ label: 'Course:', html: courseHTML });
			const instructor = parseTag({ label: 'Instructor:', html: courseHTML });

			// Extract attendance records using regex (same pattern as Go)
			const attendanceRegex =
				/(?s)<tr>\s*<td[^>]*>(\d+)<\/td>\s*<td[^>]*>([\d/]+)<\/td>\s*<td[^>]*>\s*([a-zA-Z]+)\s*<\/td>\s*<\/tr>/g;
			const attendanceMatches = Array.from(courseHTML.matchAll(attendanceRegex));

			const records = attendanceMatches.map((m) => ({
				lecture: m[1],
				date: m[2],
				status: m[3].trim()
			}));

			results.push({
				courseName,
				instructor,
				records
			});
		}

		return results;
	} finally {
		await context.close();
	}
}

export const actions: Actions = {
	default: async ({ platform, request }) => {
		if (!platform?.env.BROWSER) {
			return {
				success: false,
				error: 'Browser rendering is not available',
				courses: null
			};
		}

		const data = await request.formData();
		const username = data.get('username') as string;
		const password = data.get('password') as string;

		if (!username || !password) {
			return {
				success: false,
				error: 'Username and password are required',
				courses: null
			};
		}

		try {
			const browser = await launch(platform.env.BROWSER);

			const courses = await scrapeAttendance(username, password, browser);

			return {
				success: true,
				courses,
				error: null
			};
		} catch (error) {
			console.error('Scraping error:', error);
			const message =
				error instanceof Error ? error.message : 'An error occurred during scraping';

			// Check for specific error types
			if (message.includes('timeout')) {
				return {
					success: false,
					error: 'Request timed out. Please try again.',
					courses: null
				};
			} else if (message.includes('navigation')) {
				return {
					success: false,
					error: 'Login failed. Please check your credentials.',
					courses: null
				};
			}

			return {
				success: false,
				error: `Scraping failed: ${message}`,
				courses: null
			};
		}
	}
};
