// accept the search form

import { plucky, search } from '$lib/form-actions';
import type { Actions } from './$types';

export const actions = {

    search: search,
    plucky: plucky
} satisfies Actions;