// Search, filter and sort a page's list in the browser. Without script the
// page shows the whole list in its own order, so this only adds to it.
//
// Rows carry what they are matched and sorted by as data attributes, and the
// choices live in the URL (?q=, genre=, length=, sort=, view=) so that a
// filtered list can be linked to or bookmarked.
(function () {
  "use strict";

  const form = document.querySelector("form.search");
  const input = document.getElementById("filter");
  const counter = document.getElementById("filter-count");
  const rows = Array.from(document.querySelectorAll("[data-search]"));
  if (!form || !input || rows.length === 0) return;
  const list = rows[0].parentElement;
  const lengthSelect = document.getElementById("length");
  const sortSelect = document.getElementById("sort");
  const genreButtons = Array.from(form.querySelectorAll("[data-genre]"));
  const viewButtons = Array.from(form.querySelectorAll("[data-view]"));

  // The page's own order is the default sort, so remember it.
  rows.forEach(function (row, index) {
    row.dataset.order = index;
  });

  function number(row, key) {
    return Number(row.dataset[key]) || 0;
  }

  function text(row, key) {
    return row.dataset[key] || "";
  }

  const byText = new Intl.Collator(undefined, { sensitivity: "base", numeric: true });
  const sorts = {
    "": function (a, b) {
      return number(a, "order") - number(b, "order");
    },
    title: function (a, b) {
      return byText.compare(text(a, "title"), text(b, "title"));
    },
    name: function (a, b) {
      return byText.compare(text(a, "name"), text(b, "name"));
    },
    longest: function (a, b) {
      return number(b, "minutes") - number(a, "minutes");
    },
    shortest: function (a, b) {
      return number(a, "minutes") - number(b, "minutes");
    },
    books: function (a, b) {
      return number(b, "count") - number(a, "count");
    },
    hours: function (a, b) {
      return number(b, "minutes") - number(a, "minutes");
    },
    // Dates are YYYY-MM-DD, so they sort as text; undated rows go last.
    newest: function (a, b) {
      return text(b, "newest").localeCompare(text(a, "newest"));
    },
  };

  // Books with no recorded length belong to no length.
  const lengths = {
    short: function (minutes) {
      return minutes > 0 && minutes < 600;
    },
    medium: function (minutes) {
      return minutes >= 600 && minutes < 1200;
    },
    long: function (minutes) {
      return minutes >= 1200;
    },
  };

  const params = new URLSearchParams(location.search);
  const state = {
    q: params.get("q") || "",
    genre: params.get("genre") || "",
    length: params.get("length") || "",
    sort: params.get("sort") || "",
    view: params.get("view") || "",
  };
  // Ignore values from the URL that this page does not offer.
  function offered(select, value) {
    return (
      select &&
      Array.from(select.options).some(function (option) {
        return option.value === value;
      })
    );
  }
  if (!offered(lengthSelect, state.length)) state.length = "";
  if (!offered(sortSelect, state.sort)) state.sort = "";
  if (
    !genreButtons.some(function (button) {
      return button.dataset.genre === state.genre;
    })
  )
    state.genre = "";
  if (
    !viewButtons.some(function (button) {
      return button.dataset.view === state.view;
    })
  )
    state.view = "";

  // Each group of buttons is a radio group: one is checked, and only the
  // checked one is in the tab order.
  function check(buttons, key, value) {
    buttons.forEach(function (button) {
      const checked = button.dataset[key] === value;
      button.setAttribute("aria-checked", String(checked));
      button.tabIndex = checked ? 0 : -1;
    });
  }

  function matches(row, terms) {
    if (
      !terms.every(function (term) {
        return row.dataset.search.includes(term);
      })
    )
      return false;
    if (state.genre && !text(row, "genres").split(" ").includes(state.genre)) return false;
    if (state.length && !lengths[state.length](number(row, "minutes"))) return false;
    return true;
  }

  function apply() {
    const terms = state.q.toLowerCase().split(/\s+/).filter(Boolean);
    let shown = 0;
    for (const row of rows) {
      const match = matches(row, terms);
      row.hidden = !match;
      if (match) shown++;
    }
    const filtered = terms.length > 0 || state.genre || state.length;
    if (!filtered) counter.textContent = "";
    else if (shown === 0) counter.textContent = "No matches";
    else counter.textContent = shown + " of " + rows.length;

    // Only reorder when the order changes; filtering just hides rows.
    if (state.sort !== sortedBy) {
      const sorted = rows.slice().sort(function (a, b) {
        return sorts[state.sort](a, b) || number(a, "order") - number(b, "order");
      });
      for (const row of sorted) list.appendChild(row);
      sortedBy = state.sort;
    }
    list.classList.toggle("list", state.view === "list");
    check(genreButtons, "genre", state.genre);
    check(viewButtons, "view", state.view);

    // Keep any other parameters and the fragment; only touch this script's own.
    const query = new URLSearchParams(location.search);
    for (const key of Object.keys(state)) {
      if (state[key]) query.set(key, state[key]);
      else query.delete(key);
    }
    const search = query.toString() ? "?" + query.toString() : "";
    if (search !== location.search) {
      history.replaceState(null, "", location.pathname + search + location.hash);
    }
  }
  let sortedBy = "";

  input.value = state.q;
  if (lengthSelect) lengthSelect.value = state.length;
  if (sortSelect) sortSelect.value = state.sort;

  input.addEventListener("input", function () {
    state.q = input.value;
    apply();
  });
  if (lengthSelect)
    lengthSelect.addEventListener("change", function () {
      state.length = lengthSelect.value;
      apply();
    });
  if (sortSelect)
    sortSelect.addEventListener("change", function () {
      state.sort = sortSelect.value;
      apply();
    });
  // Click a button to choose it; arrow keys move along the group and choose,
  // as in any radio group.
  function radioGroup(buttons, key) {
    buttons.forEach(function (button, index) {
      button.addEventListener("click", function () {
        state[key] = button.dataset[key];
        apply();
      });
      button.addEventListener("keydown", function (event) {
        const step = { ArrowRight: 1, ArrowDown: 1, ArrowLeft: -1, ArrowUp: -1 }[event.key];
        if (!step) return;
        event.preventDefault();
        const next = buttons[(index + step + buttons.length) % buttons.length];
        next.focus();
        next.click();
      });
    });
  }
  radioGroup(genreButtons, "genre");
  radioGroup(viewButtons, "view");

  form.hidden = false;
  apply();
})();
