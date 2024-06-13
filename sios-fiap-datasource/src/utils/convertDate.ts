export const convertToISO8601 = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toISOString();
}
