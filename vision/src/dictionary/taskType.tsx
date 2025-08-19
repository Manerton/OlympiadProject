export const taskTypeDict = {
  Oral: 1,
  Practic: 2,
  Testing: 3
} as const;

export const taskTypes = [
  { label: "Устный", value: taskTypeDict.Oral },
  { label: "Практичекий", value: taskTypeDict.Practic },
  { label: "Тестирование", value: taskTypeDict.Testing }
] as const;