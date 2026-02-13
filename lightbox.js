// Minimal lightbox
(function () {
  var overlay = document.createElement('div');
  overlay.className = 'lightbox';
  overlay.innerHTML = '<img class="lightbox-img" alt="">' +
    '<button class="lightbox-prev" aria-label="Previous">&#8592;</button>' +
    '<button class="lightbox-next" aria-label="Next">&#8594;</button>' +
    '<button class="lightbox-close" aria-label="Close">&#215;</button>';
  document.body.appendChild(overlay);

  var img = overlay.querySelector('.lightbox-img');
  var items = Array.from(document.querySelectorAll('.photo-grid .photo-item img'));
  var current = 0;

  if (!items.length) return;

  function show(i) {
    current = (i + items.length) % items.length;
    img.src = items[current].src;
  }

  function open(i) {
    show(i);
    overlay.classList.add('active');
    document.body.style.overflow = 'hidden';
  }

  function close() {
    overlay.classList.remove('active');
    document.body.style.overflow = '';
  }

  items.forEach(function (el, i) {
    el.style.cursor = 'pointer';
    el.addEventListener('click', function () { open(i); });
  });

  overlay.querySelector('.lightbox-close').addEventListener('click', close);
  overlay.querySelector('.lightbox-prev').addEventListener('click', function (e) { e.stopPropagation(); show(current - 1); });
  overlay.querySelector('.lightbox-next').addEventListener('click', function (e) { e.stopPropagation(); show(current + 1); });
  overlay.addEventListener('click', function (e) { if (e.target === overlay) close(); });

  document.addEventListener('keydown', function (e) {
    if (!overlay.classList.contains('active')) return;
    if (e.key === 'Escape') close();
    if (e.key === 'ArrowLeft') show(current - 1);
    if (e.key === 'ArrowRight') show(current + 1);
  });
})();
