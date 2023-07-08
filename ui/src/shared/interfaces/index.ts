export interface FormButtons {
  label: string;
  class?: string;
  disableInvalidForm?: boolean;
  disabled?: boolean;
  type: "cancel" | "primary" | "secondary";
}
export interface TableColumns {
  label: string;
  propertyKey: string;
  getValueFunction?: Function;
  type?: 'progress'|'date';
}
