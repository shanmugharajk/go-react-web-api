const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

const isValidUUID = (uuid) => {
  return uuidRegex.test(uuid);
}

module.exports = {
  uuidRegex,
  isValidUUID
}