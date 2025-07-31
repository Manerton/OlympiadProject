<?php

namespace App\Components;

class DynamicFormWidget
{
    public $attributes;
    public function render()
    {
        return view('layouts.dynamic-form-widget', [
            'attributes' => $this->attributes
        ]);
    }
}
