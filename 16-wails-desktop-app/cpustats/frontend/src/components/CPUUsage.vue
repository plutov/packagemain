<template>
  <apexchart type="radialBar" :options="options" :series="series"></apexchart>
</template>

<script>
export default {
  data() {
    return {
      series: [0],
      options: {
        labels: ['CPU Usage']
      }
    };
  },
  mounted: function() {
    wails.Events.On("cpu_usage", cpu_usage => {
      if (cpu_usage) {
        this.series = [ cpu_usage.avg ];
      }
    });
  }
};
</script>
