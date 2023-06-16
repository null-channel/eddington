import { UiNodeInputAttributes, UiNodeMeta } from "@ory/client";

function getValidation(attributes: UiNodeInputAttributes) {
  const validationList = [];
  if (attributes.required) validationList.push("required");
  if (attributes.type == "email") validationList.push(attributes.type);
  return validationList.join("|");
}
const oryNodeToFormKitNodeMapper = (
  meta: UiNodeMeta,
  attributes: UiNodeInputAttributes,
  id: string
) => {
  const formElement = {
    $formkit: attributes.type,
    name: attributes.name,
    label: attributes.label,
    disabled: attributes.disabled,
    value: attributes.value,
    validation: getValidation(attributes),
  };
  if (attributes.type=='submit')
  return formElement;
};
export default oryNodeToFormKitNodeMapper;
