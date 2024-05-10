# crazy-firewall

Запрещено все что не разрешено! Закрываем доступ ко всем сайтам

##  Условия работы

- [*] Конфиг лежит на GitHub и доступен по сети
- [*] Через интервал времени получение значений сонфига, меняет плученные по необходимости
- [*] Приводит плученные значение в нижний регистр 
- [*] Сравнивает пользователя со списком пользователей в конфиге в нижнинем регистре
- [*] Сравнивает сайт обращения с разрешенным списком, если нашел сайт откроется

## Установка прокси

    Поиск -> Свойство браузера -> Подключениея -> Настройка сети -> Прокси сервер
    Или
    Сеть и Интернет -> Прокси сервер -> Настройка прокси вручную

Параметры 127.0.0.1 порт 20000

## Запрет изменения прокси

Шаг 1: Открытие редактора групповых политик

    Нажмите клавиши Win + R, чтобы открыть окно "Выполнить".
    Введите gpedit.msc и нажмите Enter. Это откроет редактор локальной групповой политики.

Шаг 2: Навигация к настройкам прокси В левой панели редактора групповых политик перейдите к:

    Пользовательская конфигурация > Административные шаблоны > Компоненты Windows > Internet Explorer.
    Найдите и выберите пункт Запретить изменение настроек прокси.

Чтобы создать свою собственную службу в Windows, которая запускается при старте системы и автоматически перезапускается при сбое, вам нужно выполнить несколько шагов.

Шаг 1: Использование sc.exe для создания службы

    Откройте командную строку от имени администратора:
        Нажмите Пуск, введите cmd, затем кликните правой кнопкой мыши на "Командная строка" и выберите "Запустить от имени администратора".

    Создайте новую службу с помощью sc.exe:
        sc create Название_Службы binPath= "Путь\к\вашему\файлу.exe" start= auto

        Название_Службы - это имя, которое вы хотите присвоить вашей службе.
        binPath - путь к исполняемому файлу вашего приложения.
        start= auto означает, что служба будет автоматически запускаться при старте системы.

Шаг 2: Настройте службу на автоматический перезапуск при сбое:
    
    Перейдите в "Управление компьютером" (правой кнопкой мыши по "Этот компьютер" > "Управление").
    В "Службы и приложения" > "Службы" найдите вашу службу.
    Кликните правой кнопкой мыши по службе и выберите "Свойства".
    Перейдите на вкладку "Восстановление".
    Настройте действия для "Первая неудача", "Вторая неудача" и "Последующие неудачи" на "Перезапустить службу".
    Укажите задержку перед перезапуском, если это необходимо.

Шаг 3: Запуск службы

    Вы можете запустить службу через командную строку или через "Управление компьютером":
    sc start Название_Службы