document.addEventListener("DOMContentLoaded", () => {
  const featureForm = document.getElementById("feature-form");
  const bugForm = document.getElementById("report-form");

  if (featureForm) {
    featureForm.addEventListener("submit", async (event) => {
      event.preventDefault();

      const formData = new FormData(featureForm);
      const data = Object.fromEntries(formData.entries());

      const response = await fetch("/api/feature", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      featureForm.reset();
    });
  }

  if (bugForm) {
    bugForm.addEventListener("submit", async (event) => {
      event.preventDefault();

      const formData = new FormData(featureForm);
      const data = Object.fromEntries(formData.entries());

      const response = await fetch("/api/report", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      bugForm.reset();
    });
  }
});
