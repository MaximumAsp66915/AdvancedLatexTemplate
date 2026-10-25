# تمرین (Homework)

ساخت خروجی از داخل پوشه‌ی `latex/`:

```bash
latexmk -xelatex main.tex
```

| مسیر | کاربرد |
|------|--------|
| `latex/main.tex` | نقطه‌ی شروع. پیش‌فرض‌های قالب در نخستین بلوک `\atlset{...}` است |
| `latex/config.tex` | عنوان، نام، درس، پالت رنگ و پوسته‌ی کد شما (به‌وسیله‌ی `mktex` نوشته می‌شود) |
| `latex/sections/` | یک پرونده برای هر مسئله. از محیط‌های `problem` و `solution` استفاده کنید |
| `latex/appendix/` | کدها |
| `latex/fig/` | شکل‌ها (بدون نام پوشه ارجاع داده می‌شوند) |
| `latex/preamble/` | `packages.tex` (بسته‌های اضافه) و `local.tex` (ماکروهای شما) |
| `latex/atl/` و `latex/fonts/` | موتور قالب و فونت‌ها که `mktex` کپی می‌کند؛ معمولاً ویرایش نمی‌شوند |
| `preprovided/`، `simulation/` | مواد درس و پرونده‌های شبیه‌سازی |

برای تغییر رنگ‌ها `palette = ...` را در `config.tex` عوض کنید
(blue، red، gray، cyan، magenta، yellow، green، sunset، ocean، berry، neon).

نکته‌ها: کلمه‌های لاتین را با `\lr{...}` بنویسید و برای نیم‌فاصله از کاراکتر ZWNJ استفاده کنید
(در ویرایشگر معمولاً `Shift+Space` یا افزونه‌ی فارسی‌نویسی).
