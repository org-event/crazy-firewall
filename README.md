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

Чтобы создать свою собственную службу в Windows, которая запускается при старте системы и автоматически перезапускается при сбое, вам нужно выполнить несколько шагов. Вот общий план действий:
Шаг 2: Использование sc.exe для создания службы

    Откройте командную строку от имени администратора:
        Нажмите Пуск, введите cmd, затем кликните правой кнопкой мыши на "Командная строка" и выберите "Запустить от имени администратора".

    Создайте новую службу с помощью sc.exe:
        Используйте следующую команду для создания службы: