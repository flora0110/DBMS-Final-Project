# -*- coding = utf-8 -*-
import os
import psycopg2

def animation():
    herokuCLI_command = 'heroku config:get DATABASE_URL -a anima-database-fe'
    DATABASE_URL = os.popen(herokuCLI_command).read()[:-1]

    # connect with database
    conn = psycopg2.connect(DATABASE_URL, sslmode='require')

    # create cursor
    cursor = conn.cursor()

    # set SQL create table
    SQL_create_command = '''
        CREATE TABLE IF NOT EXISTS animation (
            name VARCHAR(150),
            year INTEGER NOT NULL,
            season VARCHAR(5) NOT NULL,
            director VARCHAR(100),
            company VARCHAR(100),
            PRIMARY KEY (name, year, season)
        );
        '''
    # execute SQL
    cursor.execute(SQL_create_command)

    # open animation/animation.txt
    f = open("animation/animation.txt", "r", encoding="utf-8")
    # get the content in animation/animation.txt
    content = f.readlines()
    # process datas in content
    for datas in content:
        datas = datas.strip().split(",,")
        # init (name, year, season, director, company)
        name, year, season, director, company = datas[2], datas[0], datas[1], "NULL", "NULL"
        # get the value of director and company
        for index, data in enumerate(datas):
            if index > 2:
                if data[0] == 'D':
                    if director == "NULL":
                        director = data[2:]
                    else:
                        director += ("、" + data[2:])
                else:  # if data[0] == 'C':
                    if company == "NULL":
                        company = data[2:]
                    else:
                        company += ("、" + data[2:])
        # print out all data
        print(name, year, season, director, company)
        # set SQL insert data into table
        SQL_insert_command = f'''
            INSERT INTO animation
                (name, year, season, director, company)
                VALUES (%s, {year}, %s, %s, %s);
        '''
        # execute SQL
        cursor.execute(SQL_insert_command, (name, season, director, company))

    # close animation/animation.txt
    f.close()

    # commit change of database
    conn.commit()

    # close cursor
    cursor.close()
    # close connect
    conn.close()


def main():
    # animation() # done
    print("hello world")


if __name__ == '__main__':
    main()