import {PLOOGLE_API_KEY, PLOOGLE_API_URL, PLOOGLE_AUTHORIZATION_HEADER} from '$env/static/private';
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

    console.log(PLOOGLE_AUTHORIZATION_HEADER)
    const headers ={
        "Accept": "application/json",
        [PLOOGLE_AUTHORIZATION_HEADER]: PLOOGLE_API_KEY
    }
    console.log(headers)
    const res = await fetch(`${PLOOGLE_API_URL}/api/v2//search?q=${q}&p=${p || 1}`, {
        headers
    });

    if (res.status >= 400) {
        error(res.status, res.statusText);
    }

    const responseData = await res.json()
    if (!responseData.status || responseData.status !== "success") {
        error(500, responseData.status ?? "the server returned an empty response");
    }

    return {...responseData, showHidden: showHidden};
}

export const actions = {
    search: search,
    plucky: plucky,
    showHidden: async (event) => {
        const data = await event.request.formData();
        const target = (data.get("target") ?? "") as string;

        if (new URL(target).hostname !== event.url.hostname) {
            error(500, `Something has gone wrong`);
        }

        redirect(302, `${encodeURI(target)}&show-hidden`);
    }
} satisfies Actions;