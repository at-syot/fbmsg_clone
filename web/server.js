import path from 'node:path'
import express from 'express'
const app = express()
const port = process.env.PORT || 80

const distPath = path.resolve(process.cwd(), "dist")
app.use("/", express.static(distPath))
app.listen(port, () => {
  console.log(`running on port - ${port}`)
})
