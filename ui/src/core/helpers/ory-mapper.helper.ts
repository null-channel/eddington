import {
  UiNode,
  UiNodeAnchorAttributes,
  UiNodeImageAttributes,
  UiNodeInputAttributes,
  UiNodeMeta,
  UiNodeScriptAttributes,
  UiNodeTextAttributes,
  //  UiText,
} from "@ory/client";
import {
  isUiNodeAnchorAttributes,
  isUiNodeImageAttributes,
  isUiNodeInputAttributes,
  isUiNodeScriptAttributes,
  isUiNodeTextAttributes,
} from "@ory/integrations/ui";
function getValidation(attributes: UiNodeInputAttributes) {
  const validationList = [];
  if (attributes.required) validationList.push("required");
  if (attributes.type == "email") validationList.push(attributes.type);
  return validationList.join("|");
}
const oryInputMapper = (
  meta: UiNodeMeta,
  attributes: UiNodeInputAttributes,
  id: string
  //   messages:UiText
) => {
  switch (attributes.type) {
    case "hidden":
      return {
        $formkit: attributes.type,
        name: attributes.name,
        placeholder: attributes.label,
        value: attributes.value,
      };
    case "text":
      return {
        $formkit: attributes.type,
        name: attributes.name,
        placeholder: meta.label?.text ?? "",
        disabled: attributes.disabled,
        value: attributes.value,
        validation: getValidation(attributes),
      };
    case "password":
      return {
        $cmp: "FormKit",
        props: {
          type: attributes.type,
          name: attributes.name,
          prefixIconClass: "!border-0 !bg-gray-200",
          "prefix-icon": "password",
          placeholder: meta.label?.text ?? "",
          validation: getValidation(attributes),
          messageClass: "flex items-start",
        },
      };

    case "email":
      return {
        $cmp: "FormKit",
        props: {
          type: attributes.type,
          name: attributes.name,
          placeholder: meta.label?.text ?? "",
          validation: getValidation(attributes),
        },
      };
    case "submit":
      return {
        $el: "div",
        children: [
          {
            $formkit: "hidden",
            name: attributes.name,
            value: attributes.value,
          },
          {
            $cmp: "FormKit",
            props: {
              type: attributes.type,
              label: meta.label?.text ?? "",
              id: `ory-id-${id}`,
              ignore: "false",
            },
          },
        ],
      };
    default:
      throw new Error(`Input of type ${attributes.type} not implemented.`);
  }
};
const oryElementMapper = (node: UiNode) => {
  if (isUiNodeScriptAttributes(node.attributes)) {
    // Add the WebAuthn script to the DOM
    // Not tested might break !!!!!!!!!!!!!!
    console.warn("Not tested might break !!!!!!!!!!!!!!");
    const attr = node.attributes as UiNodeScriptAttributes;
    const script = document.createElement("script");
    script.src = attr.src;
    script.type = attr.type;
    script.async = attr.async;
    script.referrerPolicy = attr.referrerpolicy;
    script.crossOrigin = attr.crossorigin;
    script.integrity = attr.integrity;
    document.body.appendChild(script);
    return;
  }
  if (isUiNodeAnchorAttributes(node.attributes)) {
    const attr = node.attributes as UiNodeAnchorAttributes;
    return {
      $el: "a",
      children: [attr.title],
      attrs: { href: attr.href },
    };
  }
  if (isUiNodeImageAttributes(node.attributes)) {
    const attr = node.attributes as UiNodeImageAttributes;
    return {
      $el: "img",
      attrs: { src: attr.src, width: attr.width, height: attr.height },
    };
  }
  if (isUiNodeTextAttributes(node.attributes)) {
    const attr = node.attributes as UiNodeTextAttributes;
    return {
      $el: "p",
      children: [attr.text],
    };
  }
};
const oryMapper = (node: UiNode, id: string) => {
  if (isUiNodeInputAttributes(node.attributes)) {
    return oryInputMapper(
      node.meta,
      node.attributes as UiNodeInputAttributes,
      id
      // node.messages:
      // this for validation but i will not implement it yet since we have validation running
      // and it takes too much time to customize them
    );
  }
  if (
    isUiNodeAnchorAttributes(node.attributes) ||
    isUiNodeImageAttributes(node.attributes) ||
    isUiNodeTextAttributes(node.attributes) ||
    isUiNodeScriptAttributes(node.attributes)
  ) {
    return oryElementMapper(node);
  }
};
export default oryMapper;
