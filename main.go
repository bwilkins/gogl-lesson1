package main

import (
  "log"
  "github.com/go-gl/glow/gl-core/3.3/gl"
  glfw "github.com/go-gl/glfw3"
  "io/ioutil"
  "strings"
)

type Screen struct {
  w int
  h int
}

func LoadShader(filename string, shader_type uint32) (shader uint32) {
  var shader_bytes []byte
  var shader_string string
  var shader_err error
  var status int32
  shader = gl.CreateShader(shader_type)
  shader_bytes, shader_err = ioutil.ReadFile(filename)
  shader_bytes = append(shader_bytes, []byte("\x00")[0])
  if shader_err != nil {
    log.Fatal("Could not load shader from file: ", filename)
  }

  shader_string = string(shader_bytes)
  csource := gl.Str(shader_string)
  gl.ShaderSource(shader, 1, &csource, nil)
  gl.CompileShader(shader)

  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    log.Printf("Compile error in shader %s:\n", filename)
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    l := strings.Repeat("\x00", int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(l))

    log.Fatal(l)
  }

  return
}

func LoadShaderProgram() (program uint32) {
  var vertex_shader, fragment_shader uint32


  program = gl.CreateProgram()

  vertex_shader = LoadShader("vertex_shader.txt", gl.VERTEX_SHADER)
  fragment_shader = LoadShader("fragment_shader.txt", gl.FRAGMENT_SHADER)
  gl.AttachShader(program, vertex_shader)
  gl.AttachShader(program, fragment_shader)

  gl.LinkProgram(program)

  gl.DetachShader(program, vertex_shader)
  gl.DetachShader(program, fragment_shader)

  return
}

func LoadTriangle(program uint32) (gVAO uint32) {
  var gVBO uint32
  var vertexData []float32
  var attrib_loc uint32

  gl.GenVertexArrays(1, &gVAO)
  gl.BindVertexArray(gVAO)

  gl.GenBuffers(1, &gVBO)
  gl.BindBuffer(gl.ARRAY_BUFFER, gVBO)

  vertexData = []float32{
      //    x     y     z
          0.0,  0.8,  0.0,
         -0.8, -0.8,  0.0,
          0.8, -0.8,  0.0,
  }

  gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

  attrib_loc = uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
  gl.EnableVertexAttribArray(attrib_loc)
  gl.VertexAttribPointer(attrib_loc, 3, gl.FLOAT, false, 0, nil)

  gl.BindBuffer(gl.ARRAY_BUFFER, 0)
  clearVA()

  return
}

func clearVA() {
  gl.BindVertexArray(0)
}

func Render(window *glfw.Window, program uint32, gVAO uint32) {
  gl.ClearColor(0, 0, 0, 1)
  gl.Clear(gl.COLOR_BUFFER_BIT)

  gl.UseProgram(program)
  gl.BindVertexArray(gVAO)
  gl.DrawArrays(gl.TRIANGLES, 0, 3)

  window.SwapBuffers()
  glfw.PollEvents()

  clearVA()
  gl.UseProgram(0)

}

func main() {
  glfw_init_err := glfw.Init()
  if !glfw_init_err {
    log.Fatal("GLFW: shit broke")
  }
  defer glfw.Terminate()

  screensize := Screen{800, 600}

  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 3)
  glfw.WindowHint(glfw.ContextVersionMinor, 3)
  glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
  glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
  glfw.WindowHint(glfw.OpenglDebugContext, glfw.True)

  window, glfw_openwindow_err := glfw.CreateWindow(screensize.w, screensize.h, "I can't even remember what this program does", nil, nil)
  if glfw_openwindow_err != nil {
    log.Fatal("GLFW: could not open window")
  }


  window.MakeContextCurrent()

  glew_init_err := gl.Init()
  if glew_init_err != nil {
    log.Fatal("GLEW: could not init")
  }

  log.Printf("OpenGL Version: %s", gl.GoStr(gl.GetString(gl.VERSION)))
  log.Printf("GLSL Version: %s", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
  log.Printf("Vendor: %s", gl.GoStr(gl.GetString(gl.VENDOR)))
  log.Printf("Renderer: %s", gl.GoStr(gl.GetString(gl.RENDERER)))

  program := LoadShaderProgram()

  gVAO := LoadTriangle(program)

  for !window.ShouldClose() {
    Render(window, program, gVAO)
  }

}
