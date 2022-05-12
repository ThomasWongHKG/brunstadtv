import { oldKnex } from "../oldKnex";
import { upsertLS } from "../utils";
import episodes from '../../../btv';
import { LanguageEntity } from "@/Database";

export async function createShowTranslation(p: episodes.components["schemas"]["ItemsShowsTranslations"], m, c) {
    if (m.collection != "shows_translations") {
        return
    }
    let show = (await c.database("shows").select("*").where("id", (p.shows_id as number)))[0];
    let languages = await oldKnex<LanguageEntity>("language").select("*")
    let lang = (p.languages_code as episodes.components["schemas"]["ItemsLanguages"])    
    if (lang.code != "no") {
        // We want original to be source of truth for translations until crowdin is up
        return p
    }
    console.log("updating norwegian")
    let oldLang = languages.find(l => l.CultureCode == lang.code)
    await upsertLS(oldKnex, show.legacy_title_id, oldLang, p.title)
    await upsertLS(oldKnex, show.legacy_description_id, oldLang, p.description)
    console.log("done updating norwegian")
    return p
}

export async function createSeasonTranslation(p: episodes.components["schemas"]["ItemsSeasonsTranslations"], m, c) {
    if (m.collection != "seasons_translations") {
        return
    }
    let season = (await c.database("seasons").select("*").where("id", (p.seasons_id as number)))[0];
    let languages = await oldKnex<LanguageEntity>("language").select("*")
    let lang = (p.languages_code as episodes.components["schemas"]["ItemsLanguages"])    
    if (lang.code != "no") {
        // We want original to be source of truth for translations until crowdin is up
        return p
    }
    console.log("updating norwegian")
    let oldLang = languages.find(l => l.CultureCode == lang.code)
    await upsertLS(oldKnex, season.legacy_title_id, oldLang, p.title)
    await upsertLS(oldKnex, season.legacy_description_id, oldLang, p.description)
    console.log("done updating norwegian")
    return p
}

export async function createEpisodeTranslation(p: episodes.components["schemas"]["ItemsEpisodesTranslations"], m, c) {
    if (m.collection != "episodes_translations") {
        return
    }
    let episode = (await c.database("episodes").select("*").where("id", (p.episodes_id as number)))[0];
    let languages = await oldKnex<LanguageEntity>("language").select("*")
    let lang = (p.languages_code as episodes.components["schemas"]["ItemsLanguages"])    
    if (lang.code != "no") {
        // We want original to be source of truth for translations until crowdin is up
        return p
    }
    console.log("updating norwegian")
    let oldLang = languages.find(l => l.CultureCode == lang.code)
    await upsertLS(oldKnex, episode.legacy_title_id, oldLang, p.title)
    await upsertLS(oldKnex, episode.legacy_description_id, oldLang, p.description)
    await upsertLS(oldKnex, episode.legacy_extra_description_id, oldLang, p.extra_description)
    console.log("done updating norwegian")
    return p
}