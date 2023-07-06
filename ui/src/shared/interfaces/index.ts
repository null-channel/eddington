export interface FormButtons {
  label: string;
  class?: string;
  disableInvalidForm?: boolean;
  type: "cancel" | "primary" | "secondary";
}
export interface TableColumns {
  label: string;
  propertyKey: string;
  isFiltredBy: boolean;
  getValueFunction?: Function;
}