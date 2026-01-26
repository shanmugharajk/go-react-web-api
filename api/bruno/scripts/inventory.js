const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createProductBatch = async (data) => {
  const cachedBatch = bru.getVar(data.name)
  if (cachedBatch) {
    return cachedBatch
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/inventory/batches`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create product batch: ${result.data?.error || 'Unknown error'}`)
    }

    bru.setVar(data.name, result.data.data)

    return result.data.data
  } catch (error) {
    console.error("❌ Product batch creation failed:", error.message)
    throw error
  }
}

const deleteProductBatch = async (batchId) => {
  if (!batchId) {
    console.warn("⚠️ No batch ID provided for deletion")
    return
  }

  try {
    await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/inventory/batches/${batchId}`,
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    bru.deleteVar(batchId)
  } catch (error) {
    console.error("❌ Product batch deletion failed:", error.message)
    throw error
  }
}

module.exports = {
  createProductBatch,
  deleteProductBatch
}
