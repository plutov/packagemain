import Vue from "vue";
import App from "./App.vue";
import VueApexCharts from 'vue-apexcharts'

Vue.use(VueApexCharts)
Vue.component('apexchart', VueApexCharts)

Vue.config.productionTip = false;
Vue.config.devtools = true;

import Bridge from "./wailsbridge";

Bridge.Start(() => {
  new Vue({
    render: h => h(App)
  }).$mount("#app");
});
