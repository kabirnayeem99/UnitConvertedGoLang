document.addEventListener("DOMContentLoaded", async () => {
  const baseUrl = "http://127.0.0.1:9742";
  const fromSelect = document.getElementById("from_unit");
  const toSelect = document.getElementById("to_unit");

  if (fromSelect == null || toSelect == null) return;

  const url = window.location.href;

  const pathname = url.pathname;

  if (pathname.length == 0) pathname = "length";

  try {
    const type = pathname;
    const res = await fetch(`${baseUrl}/units?type=${type}`);

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }

    const units = await res.json();

    fromSelect.innerHTML = "";
    toSelect.innerHTML = "";

    units.forEach((unit) => {
      const fromOption = new Option(unit, unit);
      const toOption = new Option(unit, unit);

      fromSelect.add(fromOption);
      toSelect.add(toOption);
    });

    const value = document.getElementById("value");

    const convertButton = document.getElementById("convert_btn");
    convertButton.addEventListener("click", async () => {
      const res = await fetch(`${baseUrl}/convert?type=${type}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          value: value.value,
          from: fromSelect.value,
          to: toSelect.value,
        }),
      });
    });
  } catch (err) {
    console.error("Failed to load units:", err);
  }
});
