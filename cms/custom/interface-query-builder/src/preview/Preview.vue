<template>
    <div>
        <v-button style="margin-bottom:10px" :loading="loading" @click="reload">Preview</v-button>
        <div>
            <div v-for="i in items">
                <span>{{i.collection}}</span>
                <span style="margin:4px">|</span>
                <span>{{i.id}}</span>
                <span style="margin:4px">|</span>
                <span>{{i.title}}</span>
            </div>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { ref } from "vue";
import { Item } from "./types";

const props = defineProps<{
    factory: () => Promise<Item[]>;
}>()

const items = ref([] as Item[])
const loading = ref(false)

const reload = async () => {
    loading.value = true;
    items.value = await props.factory();
    loading.value = false;
}
</script>
