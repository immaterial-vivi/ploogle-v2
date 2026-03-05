import adapter from '@sveltejs/adapter-auto';
import adapterNode from '@sveltejs/adapter-node';


let selectedAdapter =  !!process.env.SHOULD_USE_ADAPTER_NODE ?  adapterNode : adapter


/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: selectedAdapter(),
		csrf: {
			checkOrigin: false,
			trustedOrigins: ['localhost:3000', "ploogle.humandomestication.guide", "ploogle.netlify.app"],
		}
	},


};

export default config;
