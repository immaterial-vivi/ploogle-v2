import { PLOOGLE_API_URL, PLOOGLE_API_KEY } from '$env/static/private';
import type { PageServerLoad } from './$types';


export const csr = false;
export const load: PageServerLoad = async ({ url }) => {
    const q = url.searchParams.get('q');
    const p = url.searchParams.get('p');


    if (!q) {
        return null
    }

    const res = await fetch(`${PLOOGLE_API_URL}/search?q=${q}&p=${p || 1}`, { headers: { "Accept": "application/json", "x-ploogle-api-key": PLOOGLE_API_KEY } });
    console.log(res)

    if (res.status >= 400) {
        error(res.status, res.statusText);
    }

    const responseData= await res.json()
    if (!responseData.status || responseData.status !== "success") {
        error(500, responseData.status ?? "the server returned an empty response");
    }

    const { status, message } = responseData;

    return { status, message };
}


import type { Actions } from './$types';
import { plucky, search } from '$lib/form-actions';
import {error} from "@sveltejs/kit";

export const actions = {

    search: search,
    plucky: plucky
} satisfies Actions;