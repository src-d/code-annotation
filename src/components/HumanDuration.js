const hour = 60 * 60 * 1000;
const minute = 60 * 1000;
const second = 1000;
const units = ['h', 'm', 's'];

export default function HumanDuration({ value }) {
  let ms = value;
  const hours = Math.floor(ms / hour);
  ms -= hours * hour;
  const minutes = Math.floor(ms / minute);
  ms -= minutes * minute;
  const seconds = Math.round(ms / second);

  return [hours, minutes, seconds]
    .map((v, i) => (v > 0 ? `${v}${units[i]}` : 0))
    .filter(v => !!v)
    .join(' ');
}
