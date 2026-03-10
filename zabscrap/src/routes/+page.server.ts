import type { Actions } from './$types';
import { launch } from "@cloudflare/playwright";


export const actions: Actions = {
    default: async ({ platform, request }) => {
        if (!platform?.env.BROWSER) {
            return { success: false, message: 'Browser unavailable' };
        }

        const data = await request.formData();
        const username = data.get('username') as string;
        const password = data.get('password') as string;

        if (!username || !password) {
            return { success: false, message: 'Missing credentials' };
        }

        const browser = await launch(platform.env.BROWSER);
        const context = browser.newContext();
        const page = (await context).newPage();
    }

}
