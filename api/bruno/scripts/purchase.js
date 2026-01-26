const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createPurchaseOrder = async (data) => {
  const cachedOrder = bru.getVar(data.name)
  if (cachedOrder) {
    return cachedOrder
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/purchase-orders`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create purchase order: ${result.data?.error || 'Unknown error'}`)
    }

    bru.setVar(data.name, result.data.data)

    return result.data.data
  } catch (error) {
    console.error("❌ Purchase order creation failed:", error.message)
    throw error
  }
}

const getPurchaseOrder = async (poId) => {
  if (!poId) {
    console.warn("⚠️ No purchase order ID provided for lookup")
    return
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/purchase-orders/${poId}`,
      method: "GET",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    if (!result.data?.success) {
      throw new Error(`Failed to fetch purchase order: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Purchase order fetch failed:", error.message)
    throw error
  }
}

const updatePurchaseOrderStatus = async (poId, status) => {
  if (!poId) {
    console.warn("⚠️ No purchase order ID provided for update")
    return
  }
  if (!status) {
    console.warn("⚠️ No status provided for purchase order update")
    return
  }

  try {
    const currentOrder = await getPurchaseOrder(poId)

    const updatedItems = currentOrder.items.map((item) => ({
      id: item.id,
      productId: item.productId,
      quantityOrdered: item.quantityOrdered,
      costPrice: item.costPrice,
      sellingPrice: item.sellingPrice,
      expiresAt: item.expiresAt
    }))

    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/purchase-orders/${poId}`,
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data: {
        vendorId: currentOrder.vendorId,
        orderDate: currentOrder.orderDate,
        status,
        notes: currentOrder.notes,
        items: updatedItems
      }
    })

    if (!result.data?.success) {
      throw new Error(`Failed to update purchase order: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Purchase order update failed:", error.message)
    throw error
  }
}

const deletePurchaseOrder = async (poId) => {
  if (!poId) {
    console.warn("⚠️ No purchase order ID provided for deletion")
    return
  }

  try {
    await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/purchase-orders/${poId}`,
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    bru.deleteVar(poId)
  } catch (error) {
    console.error("❌ Purchase order deletion failed:", error.message)
    throw error
  }
}

module.exports = {
  createPurchaseOrder,
  getPurchaseOrder,
  updatePurchaseOrderStatus,
  deletePurchaseOrder
}
