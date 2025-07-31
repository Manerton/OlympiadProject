<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class InitTables extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('attendance', function (Blueprint $table) {
            $table->id();
            $table->string('application_id', 128);
            $table->integer('status');
            $table->timestamps();
            $table->index('application_id');
        });
        Schema::create('task', function (Blueprint $table) {
            $table->id();
            $table->string('event_id', 128);
            $table->string('number', 16);
            $table->timestamps();
            $table->index('event_id');
        });
        Schema::create('task_attendance', function (Blueprint $table) {
            $table->id();
            $table->foreignId('task_id')->constrained('task')->onDelete('restrict');
            $table->foreignId('attendance_id')->constrained('attendance')->onDelete('restrict');
            $table->integer('points')->default(0);
            $table->timestamps();
            $table->unique(['task_id', 'attendance_id']);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('task_attendance');
        Schema::drop('task');
        Schema::drop('attendance');
    }
}
