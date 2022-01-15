from flask import Flask,render_template,request
import os
import psycopg2


app = Flask(__name__)


@app.route('/')
def hello_world():  # put application's code here
    return render_template("animation.html")


@app.route('/raeY')
def year_page():
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
            WHERE year={year} AND season={season};
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
                WHERE season={season};
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

    end = page * 50 + 50
    prelen = len(data)
    if end > prelen:
        end = prelen
    data = data[page * 50:end]

    left = page - 1
    right = page + 1
    if page == 0:
        left = page
    if end == prelen:
        right = page

    #return render_template("year.html",datas=data, year=year, season=season, left=left, right=right, page=page)
    return render_template("year.html",datas=data, year=year, season=season, left=left, right=right, page=page)


if __name__ == '__main__':
    app.run()
    print()
