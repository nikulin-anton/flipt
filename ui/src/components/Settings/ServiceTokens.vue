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
              @click.prevent="dialogGenerateTokenVisible = true"
            >
              Generate Token
            </a>
          </div>
        </div>
        <div class="mt-6">
          <div class="media">
            <div class="media-left">
              <p class="is-size-7">read-only</p>
            </div>
            <div class="media-content">
              <div class="pl-3">
                <p class="has-text-weight-semibold">java service</p>
                <p class="is-size-7 is-italic">Last used this month</p>
              </div>
            </div>
            <div class="media-right has-text-centered">
              <button class="button is-danger is-small is-inverted">
                Delete
              </button>
            </div>
          </div>
          <div class="media">
            <div class="media-left">
              <p class="is-size-7">read/write</p>
            </div>
            <div class="media-content">
              <div class="pl-3">
                <p class="has-text-weight-semibold">go service</p>
                <p class="is-size-7 is-italic">Last used this month</p>
              </div>
            </div>
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

    <div
      id="generateTokenDialog"
      class="modal"
      tabindex="0"
      :class="{ 'is-active': dialogGenerateTokenVisible }"
      @keyup.esc="cancelGenerateToken"
    >
      <div class="modal-background" @click.prevent="cancelGenerateToken" />
      <div class="modal-content">
        <div class="container">
          <div class="columns is-centered">
            <div class="column is-two-thirds">
              <div class="box">
                <h5 class="title is-5">New Service Token</h5>
                <form>
                  <b-field label="Name">
                    <b-input v-model="newToken.name" placeholder="Name" />
                  </b-field>
                  <b-field label="Type">
                    <div class="block">
                      <b-radio
                        v-model="newToken.type"
                        native-value="READ_TOKEN_TYPE"
                      >
                        Read-Only
                      </b-radio>
                      <b-radio
                        v-model="newToken.type"
                        native-value="WRITE_TOKEN_TYPE"
                      >
                        Read/Write
                      </b-radio>
                    </div>
                  </b-field>
                  <hr />
                  <div class="field is-grouped">
                    <div class="control">
                      <button
                        class="button is-primary"
                        :disabled="!canGenerateToken"
                        @click.prevent="generateToken"
                      >
                        Generate Token
                      </button>
                      <button
                        class="button is-text"
                        @click.prevent="cancelGenerateToken"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                </form>
              </div>
              <button
                class="modal-close is-large"
                aria-label="close"
                @click.prevent="cancelGenerateToken"
              />
            </div>
          </div>
          "
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ServiceTokens",
  data() {
    return {
      newToken: {
        name: "",
        type: "READ_TOKEN_TYPE",
      },
      dialogGenerateTokenVisible: false,
    };
  },
  methods: {
    cancelGenerateToken() {
      this.dialogGenerateTokenVisible = false;
    },
  },
};
</script>

<style></style>
