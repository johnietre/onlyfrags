const picsDiv = document.querySelector("#pics-div");

const App = {
  data() {
    return {
      //imgPaths: ["pic1.jpeg", "pic2.jpeg", "pic3.jpeg", "pic4.jpeg"],
      frags: [],
      __unused: ""
    };
  },

  async mounted() {
    if (!(await this.getFrags())) {
      return;
    }
    this.addPics();
  },

  methods: {
    addPics() {
      let newRowDiv = () => {
        let rowDiv = document.createElement("div");
        rowDiv = document.createElement("div");
        rowDiv.style.display = "flex";
        rowDiv.style.flexDirection = "row";
        rowDiv.style.justifyContent = "space-evenly";
        rowDiv.style.alignItems = "center";
        rowDiv.style.width = "100%";
        rowDiv.style.margin = "5px 0px";
        return rowDiv;
      };
      let rowDiv = newRowDiv();
      for (var i = 0; i < this.frags.length; i++) {
        if (i % 2 == 0 && i != 0) {
          document.querySelector("#pics-div").appendChild(rowDiv);
          rowDiv = newRowDiv();
        }
        var frag = this.frags[i];
        var detailsPath = `/forsale/details?frag_id=${frag.id}`;

        var parentPath = frag.parentImg;
        var parentImg = document.createElement("img");
        parentImg.src = `${parentPath}`;
        parentImg.alt = parentPath;
        parentImg.style.width = "100%";
        parentImg.style.height = "40vh";
        var parentA = document.createElement("a");
        parentA.href = "#parent";
        parentA.style.width = "50%";
        parentA.style.height = "100%";
        parentA.appendChild(parentImg);

        var fragPath = frag.fragImg;
        var fragImg = document.createElement("img");
        fragImg.src = `${fragPath}`;
        fragImg.alt = fragPath;
        fragImg.style.width = "100%";
        fragImg.style.height = "40vh";
        var fragA = document.createElement("a");
        fragA.href = detailsPath;
        fragA.style.width = "50%";
        fragA.style.height = "100%";
        fragA.appendChild(fragImg);

        var imgContainerDiv = document.createElement("div");
        imgContainerDiv.style.display = "flex";
        imgContainerDiv.style.flexDirection = "row";
        imgContainerDiv.style.alignItems = "center";
        //imgContainerDiv.style.width = "40%";
        imgContainerDiv.style.width = "100%";
        imgContainerDiv.appendChild(parentA);
        imgContainerDiv.appendChild(fragA);

        var infoP = document.createElement("p");
        infoP.innerText = `${frag.name}\n$${frag.price}`;
        var detailsA = document.createElement("a");
        detailsA.href = detailsPath;
        detailsA.innerText = "Details";
        infoP.appendChild(document.createElement("br"));
        infoP.appendChild(detailsA);

        var itemDiv = document.createElement("div");
        itemDiv.style.width = "40%";
        itemDiv.style.textAlign = "center";
        itemDiv.appendChild(imgContainerDiv);
        itemDiv.appendChild(infoP);

        rowDiv.appendChild(itemDiv);
      }
      document.querySelector("#pics-div").appendChild(rowDiv);
    },
    // Returns true if the the operation was succussful
    async getFrags() {
      const url = new URL("/api/forsale/frags", window.location.href);
      const resp = await fetch(url);
      if (!resp.ok) {
        console.log(resp);
        return false;
      }
      const json = await resp.json();
      this.frags = json;
      return true;
    },
    __unusedFn() {
    }
  }
};

const app = Vue.createApp(App);
app.mount("#app");
