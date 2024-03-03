<template>
  <header class="header">
      <img alt="Null logo" class="logo" src="@/assets/logo.svg" width="64" height="64" />
      <div class="wrapper">
          <SignedIn>
              <SearchBar />
              <UserButton after-sign-out-url="http://localhost:5173"/>
          </SignedIn>
          <SignedOut>
              <nav>
                  <RouterLink to="/">Home</RouterLink>
                  <RouterLink to="/about">About</RouterLink>
                  <RouterLink to="/blog">Blog</RouterLink>
                  <RouterLink to="/supporters">Supporters</RouterLink>
              </nav>
              <div class ="signin">
                  <SignInButton mode="modal" />
                  <SignUpButton mode="modal" />
              </div>
          </SignedOut>
      </div>
  </header>
</template>

<script setup lang="ts">
import { SignInButton, SignedIn, SignedOut, UserButton, SignUpButton } from 'vue-clerk'
import { RouterLink } from 'vue-router'
import SearchBar from '../components/SearchBar.vue'
</script>

<style>
header {
    line-height: 1.5;
    max-height: 100vh;
    width: 100%;
    grid-column: 1 / 4;
    padding-top: 1rem;
}

.logo {
    display: block;
    margin: 0 auto 2rem;
}

nav {
    font-size: 12px;
    text-align: center;
    margin-top: 2rem;
}

button {
    display: inline-block;
    outline: none;
    cursor: pointer;
    font-weight: 600;
    border-radius: 3px;
    padding: .5rem;
    margin: .25rem;
    border: 0;
    color: #fff;
    background: #ff5000;
    line-height: 1.15;
    font-size: 16px;
    :hover {
        transition: all .1s ease;
        box-shadow: 0 0 0 0 #fff, 0 0 0 3px #1de9b6;
    }
}

nav a.router-link-exact-active {
    color: var(--color-text);
}

nav a.router-link-exact-active:hover {
    background-color: transparent;
}

nav a {
    display: inline-block;
    padding: 0 1rem;
    border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
    border: 0;
}

header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-direction: row;
    flex-wrap: wrap;
    width: 100%;
}

.signin {
    padding: 1rem;
    margin-top: 1rem;
}

@media (min-width: 1024px) {
    header {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        place-items: center;
        padding-right: calc(var(--section-gap) / 2);
        width: 100%;
    }

    .logo {
        margin: 0 2rem 0 0;
    }

    header .wrapper {
        display: flex;
        place-items: flex-start;
        justify-content: space-between;
        flex-direction: row;
        flex-wrap: wrap;
    }

    nav {
        text-align: left;
        margin-left: -1rem;
        font-size: 1rem;

        padding: 1rem 0;
        margin-top: 1rem;
    }
}
</style>
