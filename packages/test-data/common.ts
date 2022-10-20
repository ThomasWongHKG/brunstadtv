import { Directus } from '@directus/sdk';
import { BTVTypes } from './types'

export const directus = new Directus<BTVTypes>(process.env.DU_HOST, {});
