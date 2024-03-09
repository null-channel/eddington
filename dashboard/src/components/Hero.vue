<script setup>

</script>

<template>
    <div class="wrapper">
    <div class="hero">
        <p class="sd">from the darkness</p>
        <p>A new cloud rises</p>
        <h1>The Null Cloud</h1>
    </div>
    </div>
</template>

<style>
div .hero {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    margin-top: 2rem;
    padding-top: 2rem;
    width: 100%;
}

div .wrapper {
    width: 100%;
    grid-column: 1 / 4;
    height: 100%;
    background-image: url('../assets/hero-bg.svg');
    background-size: cover;
}

h1 {
    font-size: 10rem;
    color: silver;
    grid-column: 1/4;
    place-self: center;
}
    
.sd {
    font-size: 1rem;
    color: #30303090;
    grid-column: 2/3;
    place-self: center;
}
    
p {
    font-size: 3rem;
    color: #88888895;
    grid-column: 2/3;
    place-self: center;
}

</style>

