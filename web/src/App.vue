<script setup lang="ts"></script>

<template>
    <router-view></router-view>
</template>

<style lang="css">
/* width */
::-webkit-scrollbar {
    width: 10px;
}

/* Track */
::-webkit-scrollbar-track {
    background: unset;
    border-radius: 10px;
}

/* Handle */
::-webkit-scrollbar-thumb {
    background: #888;
    width: 5px;
    border-radius: 10px;
}

/* Handle on hover */
::-webkit-scrollbar-thumb:hover {
    background: #555;
}
</style>
