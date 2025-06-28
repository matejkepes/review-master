import type { App } from 'vue'
import VueApexCharts from 'vue3-apexcharts'

export default ({ app }: { app: App }) => {
  app.use(VueApexCharts)
}
