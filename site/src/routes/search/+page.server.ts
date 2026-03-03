import {PLOOGLE_API_KEY, PLOOGLE_API_URL} from '$env/static/private';
import type {Actions, PageServerLoad} from './$types';
import {plucky, search} from '$lib/form-actions';
import {error, redirect} from "@sveltejs/kit";

export const csr = false;
export const load: PageServerLoad = async ({url}) => {
    const q = url.searchParams.get('q');
    const p = url.searchParams.get('p');
    const showHidden = url.searchParams.has('show-hidden');

    if (!q) {
        error(400, "A search query is required to search")
    }

    const res = await fetch(`${PLOOGLE_API_URL}/search?q=${q}&p=${p || 1}`, {
        headers: {
            "Accept": "application/json",
            "x-ploogle-api-key": PLOOGLE_API_KEY
        }
    });

    if (res.status >= 400) {
        error(res.status, res.statusText);
    }

    const responseData = await res.json()
    if (!responseData.status || responseData.status !== "success") {
        error(500, responseData.status ?? "the server returned an empty response");
    }

    const {status, message} = responseData;

    console.log("beeep:", showHidden)
    return {status, message, showHidden: showHidden};
}

export const actions = {
    search: search,
    plucky: plucky,
    showHidden: async (event) => {

        const data = await event.request.formData();

        const target = encodeURI(data.get("target") ?? "");

        console.log("hello, general kenobi", target)

        redirect(302, `${target}&show-hidden`);
    }
} satisfies Actions;