import { redirect } from "@sveltejs/kit";

export const search = async (event:any) => {

    const { request } = event
    const data = await request.formData();
    const query = encodeURI(data.get("query"))

    console.log(data)
    console.log("search", query)

    redirect(303, `/search?q=${query}`);
    // console.log("search", event)
}
export const plucky = async (event:any) => {

    const { request } = event
    const data = await request.formData();
    console.log("plucky", data)

}
