const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createStockReceipt = async (data) => {
  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/stock-receipts`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create stock receipt: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Stock receipt creation failed:", error.message)
    throw error
  }
}

const getStockReceipt = async (receiptId) => {
  if (!receiptId) {
    console.warn("⚠️ No stock receipt ID provided for lookup")
    return
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/stock-receipts/${receiptId}`,
      method: "GET",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    if (!result.data?.success) {
      throw new Error(`Failed to fetch stock receipt: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Stock receipt fetch failed:", error.message)
    throw error
  }
}

module.exports = {
  createStockReceipt,
  getStockReceipt
}
