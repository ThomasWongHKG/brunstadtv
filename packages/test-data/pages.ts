import { Directus } from '@directus/sdk';
import { BTVTypes } from './types'
import { faker } from '@faker-js/faker';
import { directus } from './common';

async function genItemsForSection() : Promise<Array<any>> {
	let eps = await directus.items('episodes').readByQuery({limit: 100})
	let ep_ids : Array<number> = []
	for (let x in eps.data) {
		ep_ids.push(eps.data[x].id)
	}

	let items = []
	do {
		let ep = ep_ids[Math.floor(Math.random()*ep_ids.length)];
		items.push({
			episode_id: ep,
			type: "episode",
		})
	} while (Math.random() > 0.3)

	return items
}

async function makeCollection() : Promise<number> {
	const coll = {
		status: "published",
		items: await genItemsForSection(),
		name: faker.company.name(),
	}

	let res = await directus.items("collections").createOne(coll)

	console.log(res)

	return res.id;
}

async function makePage() : Promise<number> {
	const page = {
		code: faker.internet.domainWord(),
		status: "published",
		type: "custom",
		translations: [
			{
				title: faker.commerce.productName(),
				description: faker.hacker.phrase(),
				languages_code: "no",
			},
		],
	}

	let res = await directus.items("pages").createOne(page)
	console.log(res);
	return res.id;
}

async function makeSection(page_id : number) : Promise<void> {
	const collection_id = await makeCollection()
	const section = {
		page_id,
		type: "item",
		status: "published",
		collection_id,
		style: "featured", // TODO: Randomize
		size: "small",
		translations: [
			{
				title: faker.commerce.productName(),
				description: faker.hacker.phrase(),
				languages_code: "no",
			},
		],
		usergroups: [
			{
				usergroups_code: "bcc-members",
			}
		]
	}

	let res = await directus.items("sections").createOne(section);
	console.log(res)
}

async function makePageWithSections(num_sections : number) : Promise<void> {
	let page_id = await makePage()

	for (let i = 0; i < num_sections; i++) {
		makeSection(page_id)
	}
}

async function start() {
	await directus.auth.login({
		email: process.env.DU_ADMIN_EMAIL,
		password: process.env.DU_ADMIN_PASS,
	});

	makePageWithSections(10);
}

start();

