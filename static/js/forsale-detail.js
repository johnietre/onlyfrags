const App = {
  data() {
    return {};
  },
  methods: {
    async fragOut() {
      const query = document.location.href.split("?", 2)[1] || "";
      const searchParams = new URLSearchParams(query);
      console.log(query, searchParams.has("frag_id"));
      const url = new URL(
        `/api/forsale/frags/buy?frag_id=${searchParams.get("frag_id")}`,
        window.location.href);
      const resp = await fetch(url, {
        method: "POST",
      });
      if (!resp.ok)
        return;
      window.location.href = (
        new URL("/forsale/buy-success", document.location.href)).toString();
    }
  }
};

const app = Vue.createApp(App);
app.mount("#app");