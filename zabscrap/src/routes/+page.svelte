<script lang="ts">
	import type { PageData } from './$types';
	import { Input } from '$lib/components/ui/input/index.js';
	import Button from '$lib/components/ui/button/button.svelte';
	import Alert from '$lib/components/ui/alert/alert.svelte';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '$lib/components/ui/table/index.js';
	import type { CourseAttendance } from '$lib/models/attendance';

	let { data }: { data: PageData } = $props();
	
	let isLoading = $state(false);
	let error = $state<string | null>(null);
	let username = $state('');
	let password = $state('');
	let courses = $state<CourseAttendance[] | null>(data?.courses || null);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;
		error = null;
		courses = null;

		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		try {
			const response = await fetch('/', {
				method: 'POST',
				body: formData
			});

			const result = (await response.json()) as { 
				success: boolean
				error?: string
				courses?: CourseAttendance[]
			};

			if (!response.ok || !result.success) {
				error = result.error || 'An error occurred. Please try again.';
				return;
			}

			courses = result.courses || null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'An unexpected error occurred';
		} finally {
			isLoading = false;
		}
	}

	function handleReset() {
		courses = null;
		error = null;
		username = '';
		password = '';
	}

	const totalLectures = $derived.by(() => {
		if (!courses) return 0;
		return courses.reduce((sum: number, c: CourseAttendance) => sum + c.records.length, 0);
	});

	const totalPresent = $derived.by(() => {
		if (!courses) return 0;
		return courses.reduce(
			(sum: number, c: CourseAttendance) => 
				sum + c.records.filter((r) => r.status.toLowerCase() === 'present').length,
			0
		);
	});

	const attendancePercentage = $derived.by(() => {
		return totalLectures > 0 ? Math.round((totalPresent / totalLectures) * 100) : 0;
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 p-4">
	<div class="mx-auto max-w-4xl">
		{#if !courses}
			<!-- Form View -->
			<div class="flex min-h-screen items-center justify-center">
				<div class="w-full max-w-md rounded-lg border border-slate-700 bg-slate-800 p-8 shadow-xl">
					<div class="mb-6 text-center">
						<h1 class="text-2xl font-bold text-white">ZabDesk Attendance</h1>
						<p class="mt-2 text-sm text-slate-400">Get your attendance records from SZABIST</p>
					</div>

					{#if error}
						<Alert variant="destructive" class="mb-4">
							<div class="text-sm font-medium text-red-200">{error}</div>
						</Alert>
					{/if}

					<form class="space-y-4" onsubmit={handleSubmit}>
						<div class="space-y-2">
							<label for="username" class="block text-sm font-medium text-slate-200">
								Registration Number
							</label>
							<Input
								id="username"
								type="text"
								name="username"
								placeholder="e.g., 12345"
								class="bg-slate-700 text-white placeholder-slate-500"
								disabled={isLoading}
								bind:value={username}
								required
							/>
						</div>

						<div class="space-y-2">
							<label for="password" class="block text-sm font-medium text-slate-200">
								Password
							</label>
							<Input
								id="password"
								type="password"
								name="password"
								placeholder="Enter your password"
								class="bg-slate-700 text-white placeholder-slate-500"
								disabled={isLoading}
								bind:value={password}
								required
							/>
						</div>

						<Button
							type="submit"
							disabled={isLoading}
							class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600"
						>
							{#if isLoading}
								<div class="mr-2 flex items-center gap-2">
									<Spinner class="h-4 w-4" />
									<span>Scraping...</span>
								</div>
							{:else}
								Begin Scraping
							{/if}
						</Button>
					</form>

					<div class="mt-6 border-t border-slate-700 pt-4">
						<p class="text-xs text-slate-500">
							Your credentials are sent directly to SZABIST's servers and are not stored.
						</p>
					</div>
				</div>
			</div>
		{:else}
			<!-- Results View -->
			<div>
				<!-- Header -->
				<div class="mb-8 rounded-lg border border-slate-700 bg-slate-800 p-6">
					<div class="flex items-center justify-between">
						<div>
							<h1 class="text-3xl font-bold text-white">Attendance Results</h1>
							<p class="mt-2 text-slate-400">Your SZABIST attendance records</p>
						</div>
						<button
							onclick={handleReset}
							class="rounded bg-slate-700 px-4 py-2 text-slate-200 hover:bg-slate-600"
						>
							← Check Another Student
						</button>
					</div>
				</div>

				{#if courses.length === 0}
					<Alert class="mb-6">
						<div class="text-sm font-medium text-slate-200">No courses found.</div>
					</Alert>
				{:else}
					<!-- Summary Stats -->
					<div class="mb-8 grid gap-4 md:grid-cols-3">
						<div class="rounded-lg border border-slate-700 bg-slate-800 p-4">
							<div class="text-sm text-slate-400">Total Lectures</div>
							<div class="mt-1 text-2xl font-bold text-white">{totalLectures}</div>
						</div>
						<div class="rounded-lg border border-slate-700 bg-slate-800 p-4">
							<div class="text-sm text-slate-400">Present</div>
							<div class="mt-1 text-2xl font-bold text-green-400">{totalPresent}</div>
						</div>
						<div class="rounded-lg border border-slate-700 bg-slate-800 p-4">
							<div class="text-sm text-slate-400">Attendance %</div>
							<div class="mt-1 text-2xl font-bold text-blue-400">{attendancePercentage}%</div>
						</div>
					</div>

					<!-- Courses -->
					<div class="space-y-6">
						{#each courses as course (course.courseName)}
							{@const coursePresent = course.records.filter((r) => r.status.toLowerCase() === 'present').length}
							{@const courseAbsent = course.records.length - coursePresent}
							{@const coursePercentage = course.records.length > 0 ? Math.round((coursePresent / course.records.length) * 100) : 0}
							<div class="rounded-lg border border-slate-700 bg-slate-800 p-6">
								<div class="mb-4">
									<h2 class="text-xl font-bold text-white">{course.courseName}</h2>
									<p class="text-sm text-slate-400">Instructor: {course.instructor}</p>
								</div>

								<!-- Attendance Table -->
								<div class="overflow-x-auto">
									<Table>
										<TableHeader>
											<TableRow>
												<TableHead class="text-slate-300">Lecture</TableHead>
												<TableHead class="text-slate-300">Date</TableHead>
												<TableHead class="text-slate-300">Status</TableHead>
											</TableRow>
										</TableHeader>
										<TableBody>
											{#each course.records as record (record.lecture)}
												<TableRow class="border-slate-700 hover:bg-slate-700/50">
													<TableCell class="text-slate-200">{record.lecture}</TableCell>
													<TableCell class="text-slate-200">{record.date}</TableCell>
													<TableCell>
														<span
															class={`inline-block rounded px-2 py-1 text-xs font-semibold ${
																record.status.toLowerCase() === 'present'
																	? 'bg-green-500/20 text-green-300'
																	: 'bg-red-500/20 text-red-300'
															}`}
														>
															{record.status}
														</span>
													</TableCell>
												</TableRow>
											{/each}
										</TableBody>
									</Table>
								</div>

								<!-- Course Stats -->
								<div class="mt-4 flex gap-4 text-sm">
									<div class="text-slate-400">
										<span class="text-green-400">{coursePresent}</span> Present
									</div>
									<div class="text-slate-400">
										<span class="text-red-400">{courseAbsent}</span> Absent
									</div>
									<div class="text-slate-400">
										<span class="text-blue-400">{coursePercentage}%</span> Attendance
									</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>
