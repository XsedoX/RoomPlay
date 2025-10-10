export default interface ISettingsSelectProps {
  value: string;
  onUpdate: (newValue: string) => void;
  items: string[];
}
