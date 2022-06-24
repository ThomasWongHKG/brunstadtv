import {SourceFiles, SourceStrings } from "@crowdin/crowdin-api-client";
import {Event, Model} from ".";
import { getConfig, getCredentials } from "./config";
import { createFile } from "./file";

type LanguagesCode = {
    code: string;
}

type Translation = {
    title: string;
    description: string;
    languages_code: LanguagesCode;
}

type TranslationPayloads = {
    show: Translation & {shows_id: number;};
    season: Translation & {seasons_id: number;};
    episode: Translation & {episodes_id: number; extra_description: string;};
}

type TranslationPayload<S extends Model> = TranslationPayloads[S]

export function getTranslationsFromEvent(input: Event<any>) {
    const payload = input.payload as Translation;
    let collection: string;
    let id: number;

    const language = payload.languages_code.code
    const values: {
        [key: string]: string;
    } = {}

    if (payload.title)
        values.title = payload.title;
    if (payload.description)
        values.description = payload.description;

    switch (input.collection) {
        case "shows_translations":
            collection = "shows"
            id = (payload as TranslationPayload<"show">).shows_id
            break;
        case "seasons_translations":
            collection = "seasons"
            id = (payload as TranslationPayload<"season">).seasons_id
            break;
        case "episodes_translations":
            collection = "episodes"
            id = (payload as TranslationPayload<"episode">).episodes_id
            let extraDescription = (payload as TranslationPayload<"episode">).extra_description
            if (extraDescription)
                values.extra_description = extraDescription
            break;
    }

    return {
        collection,
        id,
        values,
        language,
    }
}

const fileIdByCollection: {
    [collection: string]: number
} = {}

export async function getFileIdForModel(collection: string): Promise<number | null> {
    if (fileIdByCollection[collection]) {
        return fileIdByCollection[collection]
    }

    console.log("Retrieving fileId from crowdin.")
    const config = getConfig()
    const fileApi = new SourceFiles(getCredentials())
    const files = await fileApi.listProjectFiles(config.projectId, {
        directoryId: config.directoryId
    })

    for (const file of files.data) {
        if (file.data.title === collection) {
            return fileIdByCollection[collection] = file.data.id
        }
    }

    return null
}

export async function updateOrSetTranslationAsync(input: Event<any>) {
    if (input.collection.endsWith("_translations")) {
        const config = getConfig();
        const {collection, id, values} = getTranslationsFromEvent(input)

        if (!collection)
            return;

        let fileId = await getFileIdForModel(collection)

        if (fileId) {
            const stringApi = new SourceStrings(getCredentials())
            const strings = await stringApi.listProjectStrings(config.projectId, {
                fileId: fileId ?? undefined,
            })
            for (const [field, value] of Object.entries(values)) {
                const identifier = `${collection}-${id}-` + field;
                const existingString = strings.data.find(s => s.data.identifier === identifier)
    
                if (existingString) {
                    await stringApi.editString(config.projectId, existingString.data.id, [
                        {
                            op: "replace",
                            path: "/text",
                            value,
                        }
                    ])
                } else {
                    if (fileId) {
                        await stringApi.addString(config.projectId, {
                            fileId: fileId,
                            text: value,
                            identifier: `${collection}-${id}-` + field,
                        })
                    }
                }
            }
        } else {
            await createFile(collection, Object.entries(values).map(([field, value]) => ({
                identifier: `${collection}-${id}-` + field,
                context: "",
                sourceString: value
            })))
        }
    }
}
