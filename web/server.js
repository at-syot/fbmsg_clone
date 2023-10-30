import path from 'node:path'
import express from 'express'
const app = express()
const port = process.env.PORT || 80

const distPath = path.resolve(process.cwd(), "dist")
app.use("/", express.static(distPath))
app.use("/clear-data", (_, res) => {
  console.log(process.cwd() + "/cleardata.html")
  res.sendFile(process.cwd() + "/cleardata.html")
})
app.listen(port, () => {
  console.log(`running on port - ${port}`)
})
