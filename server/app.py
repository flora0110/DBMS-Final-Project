from flask import Flask,render_template,request
import os
import psycopg2


app = Flask(__name__)
app._static_folder = './static'

@app.route('/')
def hello_world():  # put application's code here
    return anima_page()


@app.route('/amina')
def anima_page():
    # set the parameter of GET
    page = request.args.get('page', default=0, type=int)
    year = request.args.get('year', default="all")
    season = request.args.get('season', default="all")

    herokuCLI_command = 'heroku config:get DATABASE_URL -a anima-database-fe'
    DATABASE_URL = os.popen(herokuCLI_command).read()[:-1]

    # connect with database
    conn = psycopg2.connect(DATABASE_URL, sslmode='require')

    # create cursor
    cursor = conn.cursor()


    # set SQL select table
    SQL_select_command = """
        SELECT * FROM animation
    """
    if year != "all" and season != "all":
        SQL_select_command = f"""
            SELECT * FROM animation
            WHERE year={year} AND season='{season}';
        """
    else:
        if year != "all":
            SQL_select_command = f"""
                SELECT * FROM animation
                WHERE year={year};
            """
        elif season != "all":
            SQL_select_command = f"""
                SELECT * FROM animation
                WHERE season='{season}';
            """

    # execute SQL
    cursor.execute(SQL_select_command)
    data = cursor.fetchall()
    # commit change of database
    conn.commit()
    # close cursor
    cursor.close()
    # close connect
    conn.close()

    # get 50 data of less
    end = page * 50 + 50
    prelen = len(data)
    if end > prelen:
        end = prelen
    data = data[page * 50:end]

    # set the page information
    left = page - 1
    right = page + 1
    if page == 0:
        left = page
    if end == prelen:
        right = page

    return render_template("anima.html",datas=data, year=year, season=season, left=left, right=right, page=page)



@app.route('/retcarahC')
def character_page():
    # set the parameter of GET
    page = request.args.get('page', default=0, type=int)
    anima_name = request.args.get('anima_name', default="all")

    herokuCLI_command = 'heroku config:get DATABASE_URL -a anima-database-fe'
    DATABASE_URL = os.popen(herokuCLI_command).read()[:-1]

    # connect with database
    conn = psycopg2.connect(DATABASE_URL, sslmode='require')

    # create cursor
    cursor = conn.cursor()

    ####################### here{

    # set SQL select table
    SQL_select_command = """
        SELECT * FROM character
    """
    if anima_name != "all":
        SQL_select_command = f"""
                SELECT * FROM character WHERE anima like '%{ anima_name }%';
        """

    ####################### }here

    # execute SQL
    cursor.execute(SQL_select_command)
    data = cursor.fetchall()

    # get 50 data of less
    end = page * 50 + 50
    prelen = len(data)
    if end > prelen:
        end = prelen
    data = data[page * 50:end]

    # execute SQL
    cursor.execute(SQL_select_command)
    # close cursor
    cursor.close()
    # close connect
    conn.close()

    # set the page information
    left = page - 1
    right = page + 1
    if page == 0:
        left = page
    if end == prelen:
        right = page

    return render_template("character.html", datas=data, anima_name=anima_name, left=left, right=right, page=page)


@app.route('/ecioV')
def voice_page():
    # set the parameter of GET
    page = request.args.get('page', default=0, type=int)
    voice_name = request.args.get('voice_name', default="all")

    herokuCLI_command = 'heroku config:get DATABASE_URL -a anima-database-fe'
    DATABASE_URL = os.popen(herokuCLI_command).read()[:-1]

    # connect with database
    conn = psycopg2.connect(DATABASE_URL, sslmode='require')

    # create cursor
    cursor = conn.cursor()

    # set SQL select table
    SQL_select_command="""
        SELECT * FROM voice
    """
    if voice_name != "all":
        SQL_select_command = f"""
                SELECT * FROM voice WHERE name like '%{voice_name}%';
        """


    # execute SQL
    cursor.execute(SQL_select_command)
    data = cursor.fetchall()

    # get 50 data of less
    end = page * 50 + 50
    prelen = len(data)
    if end > prelen:
        end = prelen
    data = data[page * 50:end]

    # get the character and the anima which voice played
    for i, a_data in enumerate(data):
        SQL_select_command = f"""
            SELECT name, anima FROM character WHERE voice like '%{a_data[0]}%';
        """
        # execute SQL
        cursor.execute(SQL_select_command)
        details_data = cursor.fetchall()

        # add detail to the tail of data
        data[i] = list(data[i])
        data[i].append(details_data)

    # close cursor
    cursor.close()
    # close connect
    conn.close()

    # set the page information
    left = page - 1
    right = page + 1
    if page == 0:
        left = page
    if end == prelen:
        right = page

    return render_template("voice.html", datas=data, character_name=voice_name, left=left, right=right, page=page)



if __name__ == '__main__':
    app.run()
    print()
