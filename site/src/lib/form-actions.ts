import { redirect } from "@sveltejs/kit";
import { PLOOGLE_API_URL, PLOOGLE_API_KEY } from '$env/static/private';
import { error } from '@sveltejs/kit';

export const search = async (event:any) => {

    const { request } = event
    const data = await request.formData();
    const query = encodeURI(data.get("query"))

    console.log(data)
    console.log("search", query)

    if (!query) {
        return
    }

    redirect(303, `/search?q=${query}`);
    // console.log("search", event)
}
export const plucky = async (event:any) => {

    const { request } = event
    const data = await request.formData();

    const query = encodeURI(data.get("query"))


    const res = await fetch(`${PLOOGLE_API_URL}/api/v2/plucky?q=${query}`, { headers: { "Accept": "application/json", "x-ploogle-api-key": PLOOGLE_API_KEY } });

    if (res.status >= 400) {
        error(res.status, res.statusText);
    }

    const responseData= await res.json()
    if (!responseData.status || responseData.status !== "success") {
        error(500, responseData.status ?? "the server returned an empty response");
    }
    console.log("plucky link", responseData)

    redirect(303, responseData.message);
}
