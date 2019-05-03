<template>
  <div id="app">
    <CPUUsage/>
  </div>
</template>

<script>
import CPUUsage from "./components/CPUUsage.vue";
import "./assets/css/main.css";

export default {
  name: "app",
  components: {
    CPUUsage
  }
};
</script>
