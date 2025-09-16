<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     *
     * @return void
     */
    public function run()
    {
        $this->call([
            AttendanceSeeder::class,
            EventScoreSeeder::class,
            TaskSeeder::class,
            TaskAttendanceSeeder::class,
            AppealSeeder::class,
        ]);
    }
}
