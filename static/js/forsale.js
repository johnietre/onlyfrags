const picsDiv = document.querySelector("#pics-div");

const App = {
  data() {
    return {
      imgPaths: ["pic1.jpeg", "pic2.jpeg", "pic3.jpeg",  "pic4.jpeg"],
      __unused: ""
    };
  },

  mounted() { 
    this.addPics();
  },

  methods: {
    addPics() {
      document.body.appendChild(document.createElement("test-comp"));
      for (var i = 0; i < this.imgPaths.length; i += 2) {
        var div = document.createElement("div");
        div.style.display = "flex";
        div.style.flexDirection = "row";
        div.style.justifyContent = "space-evenly";
        div.style.alignItems = "center";
        div.style.width = "100%";
        div.style.margin = "5px 0px";

        var name = this.imgPaths[i];
        var img = document.createElement("img");
        img.src = `/static/test-images/${name}`;
        img.alt = name;
        img.style.width = "40%";
        img.style.height = "40vh";
        var a = document.createElement("a");
        a.href = "#";
        a.style.width = "100%";
        a.style.height = "100%";
        a.appendChild(img);
        div.appendChild(a);

        if (i + 1 != this.imgPaths.length) {
          name = this.imgPaths[i + 1];
          img = document.createElement("img");
          img.src = `/static/test-images/${name}`;
          img.alt = name;
          img.style.width = "40%";
          img.style.height = "40vh";
          div.appendChild(img);
        }

        document.querySelector("#pics-div").appendChild(div);
      }
    },
    __unusedFn() {
    }
  }
};

const app = Vue.createApp(App);

app.component("test-comp", {
  template: `<h1>Yes</h1>`
});

app.mount("#app");
