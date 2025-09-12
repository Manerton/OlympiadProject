const formatDateForInput = (dateStr: string) => {
    if (!dateStr) return "";
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return ""; // если дата некорректна
    return date.toISOString().split("T")[0]; // берём только YYYY-MM-DD
};

export default formatDateForInput;