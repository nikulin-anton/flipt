<template>
  <div>
    <section class="section">
      <div class="container">
        <div class="columns">
          <div class="column is-two-thirds">
            <h5 class="title is-5">Service Tokens</h5>
            <p class="subtitle is-7">
              Tokens you have generated that allow your client code to
              authenticate with the Flipt API
            </p>
          </div>
          <div class="column">
            <a
              class="button is-primary is-normal is-pulled-right"
              @click="generateToken"
            >
              Generate Token
            </a>
          </div>
        </div>
        <div class="mt-6">
          <div class="media">
            <div class="media-left">
              <p class="has-text-weight-semibold">java service</p>
              <p class="is-size-7 is-italic">Last used this month</p>
            </div>
            <div class="media-content"></div>
            <div class="media-right has-text-centered">
              <button class="button is-danger is-small is-inverted">
                Delete
              </button>
            </div>
          </div>
          <div class="media">
            <div class="media-left">
              <p class="has-text-weight-semibold">go service</p>
              <p class="is-size-7 is-italic">Last used this month</p>
            </div>
            <div class="media-content"></div>
            <div class="media-right has-text-centered">
              <button class="button is-danger is-small is-inverted">
                Delete
              </button>
              <br />
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
export default {};
</script>

<style></style>
