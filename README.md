<h1 align='center'>Проект "Распределенный вычислитель арифметических выражений"</h1>
<br></br>
<h2>Техническое описание проекта:</h2>
<h3>Проект содержит каталоги:</h3>
<h4><b>1)</b>calculate.go - главный файл программы</h4>
<h4><b>2)</b>home.html - шаблон главной страницы</h4>
<br></br>
<h2>Инструкция по настройке и запуску проекта</h2>
<h4>Перед тем как запускать файлы, скачайте их в один проект. Для настройки и запуска проекта необходимо создать виртуальное окружение.</h4>
<br></br>
<h2>Реализация проекта через командную строку:</h2>
<h4><b>1)</b>Перед началом откройте Visual Studio Code и запустите без отладки файл calculate.go</h4>
<h4><b>2)</b>В командной строке наберите: <ul>
  <li>curl -X GET "http://localhost:8080//api/v1/calculate?exp=8%2B9" (Если выражение идет со знаком +)</li>
  <li>curl -X GET "http://localhost:8080//api/v1/calculate?exp=8%2D9" (Если выражение идет со знаком -)</li>
  <li>curl -X GET "http://localhost:8080//api/v1/calculate?exp=8%2A9" (Если выражение идет со знаком *)</li>
  <li>curl -X GET "http://localhost:8080//api/v1/calculate?exp=8%2F9" (Если выражение идет со знаком /)</li>
</ul>
</h4>
<h4>И тогда программа занесет выражение в файл, а пользователю выдаст ID его выражения:</h4>
<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/tT23h4zn/1.jpg' border='0' alt='1'/></a>
<h4><b>3)</b>В командной строке наберите curl -X GET "http://localhost:8080//api/v1/expressions". И программа выдаст список словарей - всех выражений, которые были заданы программе:</h4>
<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/9ffZch2W/c1bf1cfb-670f-4f50-b38f-d8cf4cddb728.jpg' border='0' alt='c1bf1cfb-670f-4f50-b38f-d8cf4cddb728'/></a>
<h4><b>4)</b>В командной строке наберите curl -X GET "http://localhost:8080//api/v1/expressions/<ID>". И программа выдаст информацию выражения по ID:</h4>
<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/cCf8F1Rc/image.png' border='0' alt='image'/></a>
<h4><b>5)</b>В командной строке наберите curl -X GET "http://localhost:8080//api/v1/internal/task/<ID>". И программа выдаст информацию выражения по ID:</h4>
<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/XYqYJ7rD/72f3cd12-0522-4136-aad0-d45461940941.jpg' border='0' alt='72f3cd12-0522-4136-aad0-d45461940941'/></a>
<br></br>
<h2>Реализация проекта через сайт:</h2>
<h4><b>1)</b>Перед началом откройте Visual Studio Code и запустите без отладки файл calculate.go</h4>
<h4><b>2)</b>Перейдите на сайт http://localhost:8080. Здесь будет форма, куда пользователь должен ввести свое выражение:</h4>
<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/qvNprbFP/9b02097c-5a6e-48f1-a2de-a1769804b439.jpg' border='0' alt='9b02097c-5a6e-48f1-a2de-a1769804b439'/></a>
<h4><b>3)</b>После ввода выражения, человека переведет на страницу, где будет инструкция дальнейших действий. При переходе на каждую ссылку, человек/пользователь будет видеть тоже самое .что и в командной строке.</h4>
